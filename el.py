# https://api.eloverblik.dk/CustomerApi/swagger/index.html
import requests
import json
from datetime import timedelta
from datetime import date


# Read refresh token from file
f = open("rt.secret", "r")
token = f.readline()

# Get data access token for subsequent requests
get_data_access_token_url = 'https://api.eloverblik.dk/CustomerApi/api/token'
headers = {
    'accept': 'application/json',
    'Authorization': 'Bearer ' + token,
}

response = requests.get(get_data_access_token_url, headers=headers)
print("get access token status code" + str(response.status_code))
data_access_token = response.json()['result']

print('---- retrieved access token ----\n')

# Get id of first meter - edit if you have more than one meter
metering_points_url = 'https://api.eloverblik.dk/CustomerApi/api/meteringpoints/meteringpoints'
headers = {
    'accept': 'application/json',
    'Authorization': 'Bearer ' + data_access_token,
}
meters = requests.get(metering_points_url, headers=headers)
if meters.status_code != 200:
    print('FAILED getting metering id')
    print(meters.status_code)
    exit()

first_meter = meters.json()['result'][0]['meteringPointId']

print('---- retrieved meter id ----\n')
print('meter id: ')
print(first_meter)

#Try to get data
dMinus4 = date.today() - timedelta(days=4)
dMinus3 = date.today() - timedelta(days=3)

meter_data = 'https://api.eloverblik.dk/CustomerApi/api/meterdata/gettimeseries/'
timeseries_data = {
    'dateFrom': str(dMinus4.year)+'-'+str(dMinus4.month)+'-'+str(dMinus4.day),
    'dateTo': str(dMinus3.year)+'-'+str(dMinus3.month)+'-'+str(dMinus3.day),
    'aggregation': 'Hour'
}
 
meter_data_url = meter_data + timeseries_data['dateFrom'] + '/' + timeseries_data['dateTo'] + '/' + timeseries_data['aggregation']
 
meter_json = {
    "meteringPoints": {
        "meteringPoint": [
            first_meter
        ]
    }
}
 
meter_data_response = requests.post(meter_data_url, headers=headers, json=meter_json)

if meter_data_response.status_code != 200:
    print("Data lookup failed")
    print(meter_data_response.status_code)
    exit()

meter_json = json.loads(meter_data_response.text)
meter_json_formatted = json.dumps(meter_json, indent=2)
print(meter_json_formatted)

## TODO post to Humio
#raw post 

humio_url_base="https://cloud.community.humio.com/"
humio_raw_path="api/v1/ingest/raw" #goto unstructured
humio_unstructured_path="/api/v1/ingest/humio-unstructured"
humio_post_url=humio_url_base+humio_raw_path#+humio_unstructured_path
humio_ingest_token= "0162b105-25d7-48d0-b1ad-cd387e6e5948"
headersHumio = {
    'accept': 'application/json',
    'Authorization': 'Bearer ' + humio_ingest_token,
}

print("calling humio with: " + meter_json_formatted)
humio_post_response = requests.post(humio_post_url, headers=headersHumio, json=meter_json_formatted)
#humio_post_response = requests.post(humio_post_url, headers=headersHumio, json="{hest}")
print(humio_post_response)

#Hvordan kan GET give mening som ingest ???

#
#curl https://cloud.community.humio.com/YOUR_HUMIO_URL/api/v1/ingest/raw \
#     -X POST \
#     -H "Authorization: Bearer $INGEST_TOKEN" \
#     -d 'My raw Message generated at "2016-06-06T12:00:00+02:00"'