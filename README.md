# ahouhpuc

Website of Ah Ouh PUC (ultimate frisbee club)

[Visit website](https://www.ahouhpuc.fr/fr.html)

## Development

You need a recent version of ruby and go installed.

```
git clone git@github.com:martoche/ahouhpuc.git
cd ahouhpuc
gem install jekyll
cd _etc/
go build
cd ..
jekyll serve
```

And visit http://localhost:4000/fr.html to see your changes

## Deployment

1. Ask someone the required credentials (ssh access)
2. Make sure your local git points to the version you want to deploy
3. run `_etc/deploy.sh`
