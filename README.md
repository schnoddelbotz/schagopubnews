# schagopubnews
Serverless Cms, Headleass Application, (written in) GO, (to) PUBlish news!

# idea

- Serverless GraphQL API implementation
- Backed by firestore
- maybe exists? :-)

# schagopubnews compile

- reads a yaml/hcl/json schema for db
- writes go code for type-safe firestore access, api endpoints / graphql "knowledge"
- writes emberjs editor views using generic/provided components
- ...
- creates custom schagopubnews serverless/cloudfunction binary / "graphql" API; docker imaged for local use; prod deploy = cfn
- creates custom schagopubnews emberjs build; docker imaged /w nginx for local use; prod deploy = bucket
- ...
- deploy to bucket
- deploy to cfn
