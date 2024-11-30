import asyncpg
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer
from sklearn.cluster import DBSCAN, AgglomerativeClustering
import numpy as np
from sklearn.metrics.pairwise import cosine_distances, cosine_similarity
from fastapi import FastAPI 

from motor.motor_asyncio import AsyncIOMotorClient
from pymongo.server_api import ServerApi

import os

from dotenv import load_dotenv
import faiss

load_dotenv()

MONGO_URL = os.environ.get('MONGO_URL') 
POSTGRES_URL = os.environ.get('POSTGRES_URL')

async def connect_to_mongo():
  client = AsyncIOMotorClient(MONGO_URL)
  db = client.fall24hack_mongo
  return db

async def connect_to_postgres():
   conn = await asyncpg.connect(POSTGRES_URL)
   return conn
   

async def question_helper(question):
   return {
      "question": int(question["question"])
   }

model = SentenceTransformer('sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2')

class OnNewQuestionIn(BaseModel):
    text: str


class OnNewAnswerIn(BaseModel):
    text: str
    user_id: int

class UsersAnn(BaseModel):
    user_id: int
    k: int
    search_type: str #codirectional opposite


app = FastAPI()

@app.on_event("startup")
async def startup():
   await connect_to_mongo()
   await connect_to_postgres()


@app.post("/on_new_question")
async def on_new_question(question: OnNewQuestionIn):
    postgres_db = await connect_to_postgres()

    # collection = db.questions
    # q = await collection.find_one()
    # del q["_id"]

    #db.question_collection.insert_one({"Кто ты из феечек винкс?": "0"})

    result = await postgres_db.fetch(query = "SELECT (text, cluster) FROM polls")

    q_vect = []
    clusters = []
    for row in result:
        row = dict(row)
        q_vect.append(model.encode(row['row'][0]))
        clusters.append(row['row'][1])
    unique_labels = np.unique(clusters)

    # Вычисляем уникальные кластеры
    cluster_means = []
    clusters = np.array(clusters)
    q_vect = np.array(q_vect)
    
    # Находим средние вектора для каждого кластера
    for label in unique_labels:
        cluster_points = q_vect[clusters == label]
        cluster_mean = cluster_points.mean(axis=0)
        cluster_means.append(cluster_mean)
    
    cluster_means = np.array(cluster_means)

    v = np.array([model.encode(question.text)])
    
    # Вычисляем косинусные расстояния
    similarities = cosine_similarity(v, cluster_means)[0]
    distances = 1 - similarities  # Косинусное расстояние = 1 - косинусное сходство
    
    print("DISTANCES: ", distances)
    # Определяем минимальное расстояние
    min_distance = distances.min()
    min_distance_label = unique_labels[distances.argmin()]

    print("MINDIST: ", min_distance, " ", min_distance_label)
    
    created_cluster = -1
    if min_distance > 0.65:  # Если расстояние больше порога, создаем новый кластер
        new_label = clusters.max() + 1
        clusters = np.append(clusters, new_label)
        print("CREATED NEW CLUSTER ", new_label)
        created_cluster = new_label
    else:  # Иначе добавляем в существующий кластер
        clusters = np.append(clusters, min_distance_label)
        print("ADDED TO EXISTED CLUSTER ", min_distance_label)
        created_cluster = min_distance_label

    print("ON EXIT CLUSTERS: ", clusters, "CLUSTERS_LEN: ", len(clusters))
    
    return {"cluster": int(created_cluster)}

@app.post("/on_new_answer")
async def on_new_answer(a: OnNewAnswerIn):
    mongo_db = await connect_to_mongo()

    a_vec = np.array([model.encode(a.text)])

    collection = mongo_db.profiles
    cluster_means = []
    clusters = []
    # Итерация по всем документам в коллекции
    async for document in collection.find():
    # Удаляем служебное поле MongoDB
        document.pop("_id", None)


        for k_cluster, v in document["clusters"].items(): #вытащили кластеры и их вектора
            clusters.append(k_cluster)
            vector = np.array(v["vector"])
            cluster_means.append(vector.mean(axis=0))

        unique_labels = np.unique(clusters)
        cluster_means = np.array(cluster_means)
        similarities = cosine_similarity(a_vec, cluster_means)[0]
        distances = 1 - similarities
        min_distance_label = unique_labels[distances.argmin()]

        new_answ = np.array([model.encode(a.text)])

        update_query = {"$set": {f"{a.user_id}.{min_distance_label}.vector" : ((new_answ * int(document[str(a.user_id)][str(min_distance_label)]["quantity"]) + a_vec) / (int(document[str(a.user_id)][str(min_distance_label)]["quantity"]) + 1)).tolist()}}

        result = await collection.update_one({}, update_query)
        
        #document[str(a.user_id)][str(min_distance_label)]["vector"] = (new_answ * int(document[str(a.user_id)][str(min_distance_label)]["quantity"]) + a_vec) / (int(document[str(a.user_id)][str(min_distance_label)]["quantity"]) + 1)



    return {"status": "ok"}

@app.post("/on_register_user")
async def on_register_user(user_id: int):
    mongo_db = await connect_to_mongo()

    postgres_db = await connect_to_postgres()
    result = await postgres_db.fetch(query = "SELECT (cluster) FROM polls")

    clusters = set()
    for row in result:
        row = dict(row)
        clusters.add(str(row["cluster"]))
        
    collection = mongo_db.profiles
    document = dict()
    document["user_id"] = user_id
    document["clusters"] = dict()
    for cluster in clusters:
        document["clusters"][cluster] = {"vector": np.array([model.encode("Aboba")]).tolist(), "quantity": 0}
    await collection.insert_one(document)
    return {"status": "ok"}

@app.post("/users_ann")
async def users_ann(cfg: UsersAnn):
    mongo_db = await connect_to_mongo()
    collection = mongo_db.profiles

    postgres_db = await connect_to_postgres()
    result = await postgres_db.fetch(query = "SELECT (cluster) FROM polls")

    clusters = set()
    for row in result:
        row = dict(row)
        clusters.add(str(row["cluster"]))

    # Извлечение документа для заданного user_id
    document = await collection.find_one({"user_id": cfg.user_id})
    
    user_vectors = []

    # Собираем вектора для заданного пользователя
    for key, value in document["clusters"].items():
        user_vectors.append(np.array(value["vector"][0]))

    # Формируем вектор для текущего пользователя
    np_user_vectors = np.array(user_vectors).astype(np.float32)  # Размерность (N, 384)
    print("User vectors shape:", np_user_vectors.shape)

    # Создаем индекс FAISS
    index = faiss.IndexFlatL2(384)  # Убедимся, что размерность совпадает

    ids = []
    other_user_vecs = []

    # Собираем вектора всех других пользователей
    async for d in collection.find():
        ids.append(d["user_id"])
        for key, value in d["clusters"].items():  # Здесь используем d, а не document
            vector = np.array(value["vector"][0])
            other_user_vecs.append(vector)

    # Преобразуем другие вектора в массив numpy
    other_user_vecs = np.array(other_user_vecs).astype(np.float32)
    print("Other user vectors shape:", other_user_vecs.shape)

    # Добавляем другие вектора в индекс
    index.add(other_user_vecs)

    # Подготовка запроса
    if cfg.search_type == "opposite":
        np_user_vectors = -np_user_vectors

    # Выполняем поиск для каждого вектора пользователя
    result = []
    for user_vector in np_user_vectors:  # Размерность (384,)
        distances, indices = index.search(user_vector.reshape(1, -1), cfg.k)
        print(distances)

        for dist, idx in zip(distances[0], indices[0]):  # Индексы для 1 вектора
            user_id = ids[idx]
            overlapping_percentage = float(dist * 100)
            result.append(dict(user_id=user_id, overlapping_percentage=overlapping_percentage))

    return result