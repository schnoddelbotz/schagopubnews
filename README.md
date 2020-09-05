# schagopubnews
Serverless Cms, Headleass Application, (written in) GO, (to) PUBlish news!

# idea

- Serverless GraphQL API implementation
- Backed by firestore
- maybe exists? :-)

# schagopubnews compile

- reads a yaml/hcl/json schema for db; maybe https://graphql.org/learn/schema/ , use https://github.com/graphql-go/graphql ??
- writes go code for type-safe firestore access, api endpoints / graphql "knowledge"
- writes emberjs editor views using generic/provided components , use apollo:
  - https://github.com/ember-graphql/ember-apollo-client
  - https://www.howtographql.com/ember-apollo/0-introduction/
  - https://medium.com/kloeckner-i/ember-and-graphql-8aa15f7a2554
- ...
- creates custom schagopubnews serverless/cloudfunction binary / "graphql" API; docker imaged for local use; prod deploy = cfn
- creates custom schagopubnews emberjs build; docker imaged /w ~~nginx~~ schagopubnews serving static content, too, for local use; prod deploy = bucket
- ...
- deploy to bucket
- deploy to cfn

# ideas

- use CloudRun for SPN producers (e.g. latex output producer ...)
- Docker image(s) for usage scenarios: plain CloudRun ... or +plugins:
- Producer / Output "plugins":
  - markdown2latex converter ... whats its name? integrate...
  - add hugo!! (support remote template source? fetch? add to image?)
  - add dumb json output ... add jq examples ... script integration of api
  - producer with support for mass- ... -e-mailing/
  - artefacts into bucket -> notification new product ready in ui :)