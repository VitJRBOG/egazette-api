# Installation

1. `git clone https://github.com/VitJRBOG/egazette-api`
2. `cd egazette-api`
3. `docker-compose up`

# Using

Copy any of the URLs `http://127.0.0.1:8000/source/[source_name]/articles` and paste them into your RSS reader.

Open any of the URL `http://127.0.0.1:8000/source/[source_name]/info` in your web browser for more information about the source.

Replace `[source_name]` with any of the [source names](#available-sources).

**P.S.**

_If you are running the application on a remote machine, you should replace `127.0.0.1` with the IP address of that remote machine._

_`8000` is the default port. You can replace it with another port in the `docker-compose.yml` file._

# Available sources

-   Jet Propulsion Laboratory News: `jpl`
-   Vestirama's News: `vestirama`
