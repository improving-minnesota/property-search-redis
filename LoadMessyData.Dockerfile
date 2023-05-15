FROM python:3.10-alpine
WORKDIR /app
COPY requirements.txt /app
RUN pip install -r requirements.txt

COPY LoadMessySampleData.py /app
# data can change often - main separate layer
COPY MessySampleData.txt /app
CMD ["python", "LoadMessySampleData.py"]