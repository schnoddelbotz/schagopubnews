# schagopubnews
Serverless Cms, Headleass Application, (written in) GO, (to) PUBlish news!

# idea

- Serverless GraphQL API implementation
- Backed by firestore
- maybe exists? :-)

# schagopubnews compile

- reads a yaml/hcl/json schema for db; maybe https://graphql.org/learn/schema/ , use https://github.com/graphql-go/graphql ??
- writes go code for type-safe firestore access, api endpoints / graphql "knowledge"
- writes emberjs editor views using generic/provided components , use https://github.com/ember-graphql/ember-apollo-client ( https://www.howtographql.com/ember-apollo/0-introduction/ )
- ...
- creates custom schagopubnews serverless/cloudfunction binary / "graphql" API; docker imaged for local use; prod deploy = cfn
- creates custom schagopubnews emberjs build; docker imaged /w nginx for local use; prod deploy = bucket
- ...
- deploy to bucket
- deploy to cfn
