# aws-get-mfa-session
A helper tool for simplifying getting a Session Token when using MFA on an AWS IAM User. Gets the MFA device details for the currently configured IAM user.

When you want a session token for an MFA protected IAM account you currently need to go dig out your MFA device serial number from the console or the CLI, and whilst they don't change much it's a bit of a annoyance. This simplifies it to entering your MFA token and getting the credentials back ready to be pasted in without transforming the CLI JSON output yourself.
