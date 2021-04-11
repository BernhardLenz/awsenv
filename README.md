# awsenv
AWS cli profile switcher.

`awsenv` lists the content of the aws credentials and config files in a better readable format and makes it easier to switch the default profile.  

## Install:
#### MacOS
Use Choco (details to follow)
#### Windows
Use Homebrew (details to follow)

## Usage:
### List all available profiles:
```sh
$ awsenv
```
or
```sh
$ awsenv list
```

Output: 
```shell
  PROFILE     AWS_ACCESS_KEY_ID       REGION       OUTPUT    
  personal    ****************LDIA    us-east-1    text          
* office      ****************NQPA    [us-east-1]  [json]      
```
Profiles with * are active profiles. 

Profiles with region or output in [] are using the default config.

### Activate a given profile:
```sh
$ awsenv activate <profile>
```
Example:
```sh
$ awsenv activate personal

$ awsenv 
  PROFILE     AWS_ACCESS_KEY_ID       REGION       OUTPUT    
* personal    ****************LDIA    [us-east-1]  [json]          
  office      ****************NQPA    us-east-1    text      
```
The activate command changes the `[default]` section in the `credentials` file. 

The activate command DOES NOT modify the `config` file.
  
### Create new profile:
Use the aws cli to create a new profile
```sh
$ aws configure
```

### Help:
```sh
$ awsenv help
```
