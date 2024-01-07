<h1 align="left">SIMPLEstack</h1>
<br/>

SIMPLEstack is a free and simplistic [content management system](https://en.wikipedia.org/wiki/Content_management_system) conceived to be used as boilerplate or starting point to create full-stack web applications using Go, HTMX and PostgreSQL.

[![Go Version](https://img.shields.io/github/go-mod/go-version/kevinsuner/SIMPLEstack)](https://github.com/kevinsuner/SIMPLEstack/blob/master/go.mod)
[![License](https://img.shields.io/github/license/kevinsuner/SIMPLEstack)](https://github.com/kevinsuner/SIMPLEstack/blob/master/LICENSE)

---

<p align="center">
    <kbd><img src="https://i.postimg.cc/wjfbnYsZ/simplestack-screenshot.png" alt="SIMPLEstack Screenshot" title="SIMPLEstack Screenshot"/></kbd>
</p>

---

## Motivation
When I started doing web development the most used stack was the [LAMP stack](https://en.wikipedia.org/wiki/LAMP_%28software_bundle%29)
which with all of it's problems was for me (at least) a very straightforward way to develop web applications. Then Javascript frameworks
started popping like mushrooms on a humid forest, it was a breadth of fresh air in the field, but the feeling that I had when using these
new tools was a lot of added complexity.

Not long after that I decided to focus solely on the back-end side of things, even though I liked (and still like) the front-end. Fast-forward to this day and apparently we've come full circle and realized that maybe??? all of this complexity wasn't needed (or even desirable) and that sending plain html to the client instead of JSON works for a lot of use cases.

After seeing this sort of back-to-basics movement that started to emerge in the tech community, I decided that I wanted to work on a full-stack web application again, just not with PHP nor Javascript but Go. Thanks to [ThePrimeagen](https://www.youtube.com/channel/UC8ENHE5xdFSwx71u3fDH5Xw) I discovered [HTMX](https://htmx.org) and from there SIMPLEstack was born and my desire to work on the front-end resurrected :) 

**NOTE**: In any way I underappreciate the effort that a lot of people way smarter than me have put in developing these new frameworks, but they just aren't for me, I prefer a more "rudimental" approach, that all. 

## Contributing
All contributions are extremely appreciated, if you find an issue that is interesting
to you, do not hesitate and say something so I know that you are hacking on that. Also
do not be afraid to ask for help or feel like you are missing some information, I will
be delighted to help :)


### Git commit message guidelines
The most important part is that each commit messages should have a title/subject in imperative
mood starting with a capital letter and no trailing period: `generator: Fix articles list being sorted in ascending order`
**NOT** `list sorted right.`. If you still unsure about how to write a good commit message 
this [blog article](https://cbea.ms/git-commit/) is a good resource for learning that.

Most title/subjects should have a lower-cased prefix with a colon and one whitespace. The prefix can be:
- The name of the package where (most of) the changes are made (e.g. `parse: Add RawTitle to metadata struct`)
- If the commit touches several packages with a common functional topic, use that as a prefix (e.g. `errors: Resolve correct line numbers`)
- If the commit touches several packages without a common functional topic, prefix with `all:` (e.g. `all: Reformat go code`)
- If this is a documentation update, prefix with (e.g. `docs:`)
- If nothing of the above applies, just leave the prefix out

Also, if your commit references some or more Github issues, always end your commit message body with *See #1234* or *Fixes #1234*.
Replace *1234* with the Github issue ID.

An example:
```text
generator: Fix articles list being sorted in ascending order

Added a function that returns a sorted "map" by passing it an unsorted one, 
it creates an "array" with the keys of the first "map", sorts the given "array"
using sort.StableSort and iterates over those keys to create new entries in the
new "map" and sets its value by accessing the first "map".

Fixes #1234
```

### Fetching the source from Github
SIMPLEstack uses the Go modules support built into Go 1.11 to build. The easiest is to clone RunGo in a directory outside of `GOPATH`,
as in the following example:
```bash
mkdir $HOME/src
cd $HOME/src
git clone https://github.com/kevinsuner/SIMPLEstack
cd SIMPLEstack 
```

Now, to make a change to SIMPLEstack's source:
1. Create a new branch for your changes (the branch name is arbitrary):
    ```bash
    git checkout -b abc123
    ```
2. After making your changes, commit them to your new branch:
    ```bash
    git commit -a -v
    ```
3. Fork RunGo in Github
4. Add your fork as a new remote (the remote name, "foo" in this example, is arbitrary):
    ```bash
    git remote add git@github.com:USERNAME/SIMPLEstack.git
    ```
5. Push the changes to your new remote
    ```bash
    git push --set-upstream foo abc123
    ```
6. You are now ready to submit a PR based upon the new branch in your forked repository

### Building SIMPLEstack with your changes
Prerequisites to build SIMPLEstack from source:
- [Go 1.21](https://go.dev/dl) or later

Run the command `go build` to build the binary. 

## Credits
SIMPLEstack makes use of a variety of open-source projects including:
- [github.com/golang/go](https://github.com/golang/go)
- [github.com/bigskysoftware/htmx](https://github.com/bigskysoftware/htmx)
- [github.com/postgres](https://github.com/postgres)
- [github.com/joho/godotenv](https://github.com/joho/godotenv)
- [github.com/lib/pq](https://github.com/lib/pq)
- [github.com/yuin/goldmark](https://github.com/yuin/goldmark)
