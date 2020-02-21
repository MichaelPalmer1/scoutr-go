# OIDC Examples

Here are examples of using SimpleAPI with either AWS or GCP in an OIDC setup.

## Assumptions
- Your go program is running behind a web server or middleware that adds the following OIDC headers:
    - Oidc-Claim-Name (Person's name)
    - Oidc-Claim-Sub (Person's username)
    - Oidc-Claim-Mail (Person's email)
    - Optional: Oidc-Claim-Groups (Comma separated list of groups the user is a member of)
