This repo provides a set of tools for playing with and testing Heroku Private Spaces VPC peering.

**WARNING**: Executing scripts in this repo will spin up AWS resources that will incur cost on your AWS account.

### Prerequisites

* AWS CLI installed with a working AWS profile. Set the profile name in the `AWS_PROFILE` environment variable.
* Heroku CLI installed and logged in. You must have Heroku Enterprise to create Private Spaces.
* `jq` installed

### Create a Private Space

You'll need a private space to play with. Don't use an existing space with important apps in it. Create a space:

	heroku spaces:create acmespace --org myorg

### Create a VPC

The following command creates a new VPC with one t2.micro EC2 instance and one db.t2.micro 
RDS Postgres instance:

    bin/createvpc acmestack 

It uses cloudformation and you pass the stack name as argument. You'll be using this stack name to execute other commands. The script will create a VPC with a 172.31.0.0/16 CIDR block. You can choose other /16 CIDR blocks by passing in an extra argument, e.g:

    bin/createvpc acmestack 172.16

to create a 172.16.0.0/16 space. The Heroku Private Space uses 10.0.0.0/16, so don't pick that CIDR block.

IMPORTANT: This script will upload your public key in $HOME/.ssh/id_rsa.pub. Things won't work if you don't have a key there and if you don't want this public key to be uploaded to your AWS account, then modify the script.

### Peer the VPC to the space

Peer the VPC to the Private Space with:

    bin/peer acmestack acmespace

This script will

1. use your AWS creds to initiate the peering connection
2. use your Heroku creds to accept the peering connection
3. use your AWS creds to set up routes between the space and the VPC

### Play with the peering

To test out the peering, this repo also happens to be a fully functional Heroku app. You can do a "button deploy" of this app which injects the EC2 instance IP address and the postgres URL for the database as config vars using the following command:

    bin/setuphello acmestack acmespace

This will deploy what's in the jesperfj/hello-peer Github repository (check the script). You can of course manually create a Heroku app in your space, push the code and set the config vars. This is just a convenience.

Once the app is deployed, you can run the very same app on the EC2 instance in the VPC. This lets you check if the peering connection is working. In a separate terminal, run:

    bin/runhello acmestack haikuappname

where haikuappname is the name of the app created by the `setuphello` script.

The `runhello` script works by grabbing the slug download URL on your machine using your Heroku credentials and then using that URL in an ssh remote executing on the EC2 instance to download the slug, extract it, and run the app inside it. It's hardcoded for the hello-peer app, but it's easy to adapt for other use cases.

### TODO

Need some samples of accessing the RDS database.
