FROM python:3.12-slim
ARG SOURCE_PATH=fastapi/asyncpg
WORKDIR /app
COPY fastapi/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY $SOURCE_PATH ./server
EXPOSE 8080
CMD ["uvicorn", "server.main:app", "--host", "0.0.0.0", "--port", "8080"]