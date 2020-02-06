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

## Install a new server

See the install.sh script. Additionnal steps:

- Copy the /home/martin/ahouhpuc/production.env file on the new server (don’t
	ever commit this file in the git repository).
- Copy the /home/martin/ahouhpuc/autocert dir on the new server (don’t
	ever commit these files in the git repository).
