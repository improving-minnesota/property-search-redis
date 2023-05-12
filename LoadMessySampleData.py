from redisearch import Client, IndexDefinition, TextField
from redis import ResponseError
import time

schema = (
    TextField("account", weight=4.0),
    TextField("owner", weight=2.0),
    TextField("address", weight=3.0),
    TextField("class", weight=1.0)
)
fld_account = 2
fld_owner = 6
fld_address = 12
fld_class = 13
create_index = False
r_client = Client("data")
start = time.time()


try:
    print(r_client.info())
except ResponseError:
    print("Index needs to be created")
    r_client.create_index(schema, definition=IndexDefinition(prefix=['doc:']))


with open('MessySampleData.txt', encoding='ISO-8859-1') as data:
    c = 0
    for line in data:
        c = c + 1
        if c == 1:
            continue        # pass the first line because it's the header
        try:
            data = line.strip().split("|")
            r_client.redis.hset("doc:" + data[fld_account], mapping={
                "account": str(data[fld_account]).rjust(8, '0'),
                "owner": data[fld_owner].strip(),
                "address": data[fld_address].strip(),
                "class": data[fld_class].strip()
            })
            if c % 1000 == 0:
                print(f'{c:,}')
        except Exception as e:
            print("Error line " + c)
            print(e)
            print("----------------------------")
            print(data)

end = time.time()
elapsed = end - start
print("Completed in %s sec" % elapsed)
