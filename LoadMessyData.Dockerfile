FROM python:3.10-alpine
WORKDIR /app
COPY requirements.txt /app
COPY LoadMessySampleData.py /app
RUN pip install -r requirements.txt
# data can change often - main separate layer
COPY MessySampleData.txt /app
CMD ["python", "LoadMessySampleData.py"]