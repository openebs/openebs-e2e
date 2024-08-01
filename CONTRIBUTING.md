# Contributing Guidelines
<BR>

## Umbrella Project
OpenEBS is an "umbrella project". Every project, repository and file in the OpenEBS organization adopts and follows the policies found in the Community repo umbrella project files.
<BR>

This project follows the [OpenEBS Contributor Guidelines](https://github.com/openebs/community/blob/HEAD/CONTRIBUTING.md)

---

## Sign your work

We use the Developer Certificate of Origin (DCO) as an additional safeguard for the OpenEBS projects. This is a well established and widely used mechanism to assure that contributors have confirmed their right to license their contribution under the project's license. Please read [dcofile](https://github.com/openebs/openebs/blob/HEAD/contribute/developer-certificate-of-origin). If you can certify it, then just add a line to every git commit message:

Please certify it by just adding a line to every git commit message. Any PR with Commits which does not have DCO Signoff will not be accepted:

````
  Signed-off-by: Random J Developer <random@developer.example.org>
````

Use your real name (sorry, no pseudonyms or anonymous contributions). The email id should match the email id provided in your GitHub profile.
If you set your `user.name` and `user.email` in git config, you can sign your commit automatically with `git commit -s`.

You can also use git [aliases](https://git-scm.com/book/tr/v2/Git-Basics-Git-Aliases) like `git config --global alias.ci 'commit -s'`. Now you can commit with `git ci` and the commit will be signed.

