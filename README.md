# Presentation:
This OSINT tool was created following the release of the following [article](https://trufflesecurity.com/blog/anyone-can-access-deleted-and-private-repo-data-github) and the rump of **Frederick Kaludis** during **Winerump 2024**, it allows you to find all the commits deleted from a GitHub repository using the API.

# Operation:
We do a first fetch on the /commits endpoint and another on the /events endpoint. We then compare the two responses by looking for events of type "PushEvent." If the event is not in the commits, then it is a deleted commit.

# Usefulness:
During a security audit or a pentest, it can be interesting to examine the commits deleted from an organization's repository to find sensitive informations.

![img](meme-leak-password.png)

**⚠️ Warning: This tool is reserved for professionals, I am not responsible if you use it illegally.**

# Features:
- Automatic recovery of all public repositories of a user
- Intuitive and simple interface

# Setup & usage

## Create a .env file

```bash
touch .env
echo 'GITHUB_TOKEN = "{your github api key}"' >> .env
```

## Build 
```bash
go build main.go
```

## Run ! 
```bash
./main
```


# Thanks

I would like to personally thank Frederic for his presentation and for taking the time to give me some additional explanations after his presentation.
