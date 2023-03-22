# Steps os Project to make a slackbot to onboarding SEBRAE-SP

## Step 01

- Construct a simple slackbot to understand how slackbot works.
- My choise was make a slackbot that calculates the age based on your year or birth.
- [See the repository here](https://github.com/alessandra1408/SlackBot_CalculeteAge)

## Step 02

- Structure what the slackbot for Onboarding needs to do:
  - Activate from a message command
  - Send a message to new user in squad
  - Send email correctly
  - Create a task to meet the operations team
  - Send a channel notification that everything is fine (if true)
    - Send notification if there was any problem in the processes

## Step 03

- Map, analyze and choose the best way to send emails with golang
- Options:
  - [Gomail](https://pkg.go.dev/gopkg.in/gomail.v2) -> It will not be used because neoway's google account does not have the required app password
  - [turnMail]()
  - [Chilkat](https://www.example-code.com/golang/smtp_noAuthentication.asp) -> They unlocked 30 days free but we need to pay after
  - [SendGrid](https://docs.sendgrid.com/for-developers/sending-email/quickstart-go) -> Needs to pass the password of email before, so it's no security and functional
