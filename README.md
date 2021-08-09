# covid19-service
There are 2 api's in this service.
1. To update the covid caches from each state and  union territories of India.
curl : 
`0.0.0.0:8000/api/v1/covid/update`

2. Second api fetches covid cases at the user's location. User insert  location in format of lat lng and receives covid cases in his state or ut and in india.
curl : 
`curl --location --request GET '0.0.0.0:8000/api/v1/covid/cases' \
--header 'Content-Type: application/json' \
--data-raw '{
    "lat":"20.2270",
    "lng":"73.0169"
}'`

Second api have in-memory cache setup of 30 min.
