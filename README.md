> [!WARNING]
> This scraper does not work anymore as SensCritique v1 ([old.senscritique.com](https://old.senscritique.com)) was shut down definitively in November 2023. The code needs to be updated for SensCritique v2, or switch to using their GraphQL API at [apollo.senscritique.com](https://apollo.senscritique.com)

# go-senscritique

A SensCritique web scraper.

## Back up user diary via SensCritique v2 GraphQL API

Set your username in `variables`. The `map` function used in the `jq` command converts the `dateDone` property from UTC to your local timezone.

```sh
$ curl -s -X "POST" "https://apollo.senscritique.com/" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -d $'{
        "query": "query UserDiary($isDiary:Boolean $limit:Int $offset:Int $universe:String $username:String! $yearDateDone:Int){user(username:$username){collection(isDiary:$isDiary limit:$limit offset:$offset universe:$universe yearDateDone:$yearDateDone){products{id universe category title originalTitle alternativeTitles yearOfProduction url otherUserInfos(username:$username){dateDone rating}}}}}",
        "variables": {
                "isDiary": true,
                "limit": 5000,
                "universe": null,
                "username": "imkh",
                "yearDateDone": null
        }
}' | jq '.data.user.collection.products | map(.otherUserInfos |= (.dateDone |= (gsub("\\.\\d+Z"; "Z") | fromdateiso8601 | strflocaltime("%Y-%m-%dT%H:%M:%S %Z"))))' > "backup_senscritique_$(date '+%Y-%m-%d').json"
```
