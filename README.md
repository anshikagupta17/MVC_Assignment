<h1 align="center">MVC_Assignment</h1>
<hr>
<h2>Description / Overview</h2>
<p>
  A Clash of Clans inspired village-building and battle game, built as a full-stack assignment project using Go (backend, MVC architecture) and   Next.js (frontend).
  <p><b>Assignment Overview</b></p>
  <p>This project was built to satisfy the following requirements:
    <ul>
    <li>Authentication: JWT-based user authentication</li>
    <li>Language / Architecture: Go, following the MVC (Model-View-Controller) pattern</li>
    <li>Database: PostgreSQL with a versioned migration system (golang-migrate)</li>
    <li>Testing: Unit tests for core logic</li>
    <li>CI/CD: A GitHub Actions pipeline that builds the project and runs tests on every push</li>
    <li>Deployment: Docker Compose for one-command setup (Postgres + backend + frontend)</li>
    <li>Data Seeding: Automatic seeding of test users, test villages, and game metadata on first run</li>
    <li>Makefile</li>
    </ul>
  </p>
</p>
<hr>
<h2>About the Game</h2>

<p>Each player starts with their own village containing 
  <ul>
    <li>Townhall</li>
    <li>Gold Mine</li>
    <li>Elixir Mine</li>
    <li>Army Camp</li>  
  </ul>
</p>
<p>Player can
  <ul>
    <li>Register a user</li>
    <li>Login their own created user or use seeded test users</li>
    <li>Build new structures (defenses, resource buildings, storage) using gold or elixir</li>
    <li>Upgrade existing buildings to increase their level, capped by their Townhall level</li>
    <li>Move buildings around their 50x50 village grid</li>
    <li>Collect resources produced over time by their mines, capped by capacity of respective storage buildings</li>
    <li>Train troops using elixir, and upgrade troop types to make them stronger</li>
    <li>Battle other players' villages — find a matched opponent, select an army, and attack to earn loot and trophies</li>
  </ul>
</p>
<hr>
<h2>Game Rules</h2>
<p>
  <ul>
    <li>Buildings cannot overlap on the grid</li>
    <li>A building's level can never exceed the village's current Townhall level</li>
    <li>Only one building can be upgrading at a time</li>
    <li>Defensive buildings (Cannon, Archer Tower, Wall, Mortar) unlock at specific Townhall levels</li>
    <li>Each building type has a maximum quantity allowed per Townhall level</li>
    <li>Resource collection is capped by total storage capacity (Gold/Elixir Storage)</li>
    <li>Troop training costs a flat elixir cost per unit; troop upgrades cost more per level and take time</li>
    <li>Matchmaking finds an opponent within ±1 Townhall level and ±100 trophies</li>
    <li>Battle outcome (stars, destruction %, loot, trophy change) is calculated from total attack power vs. total defense power</li>
  </ul>
</p>
<hr>
<h2>Building Types</h2>
<table>
  <tr>
    <th>ID</th>
    <th>Name</th>
  </tr>
  <tr>
    <td>1</td>
    <td>Townhall</td>
  </tr>
  <tr>
    <td>2</td>
    <td>Cannon</td>
  </tr>
  <tr>
    <td>3</td>
    <td>Archer Tower</td>
  </tr>
  <tr>
    <td>4</td>
    <td>Wall</td>
  </tr>
  <tr>
    <td>5</td>
    <td>Mortar</td>
  </tr>
  <tr>
    <td>6</td>
    <td>Gold mine</td>
  </tr>
  <tr>
    <td>7</td>
    <td>Elixir mine</td>
  </tr>
  <tr>
    <td>8</td>
    <td>Gold Storage</td>
  </tr>
  <tr>
    <td>9</td>
    <td>Elixir Storage</td>
  </tr>
  <tr>
    <td>10</td>
    <td>Army Camp</td>
  </tr>
</table>
<hr>
<h2>Troop Types</h2>
<table>
  <tr>
    <th>ID</th>
    <th>Name</th>
    <th>Unlock level</th>
  </tr>
  <tr>
    <td>1</td>
    <td>Barbarian</td>
    <td>1</td>
  </tr>
  <tr>
    <td>2</td>
    <td>Archer</td>
    <td>2</td>
  </tr>
  <tr>
    <td>3</td>
    <td>Giant</td>
    <td>3</td>
  </tr>
  <tr>
    <td>4</td>
    <td>Wizard</td>
    <td>4</td>
  </tr>
  <tr>
    <td>5</td>
    <td>Goblin</td>
    <td>4</td>
  </tr>
</table>
<hr>
<h2>Tech Stack</h2>
<p>
  <ul>
    <li>Backend: Go, MVC architecture</li>
    <li>Database: PostgreSQL, golang-migrate for migrations</li>
    <li>Authentication: JWT</li>
    <li>Frontend: Next.js (React+JS+CSS)</li>
    <li>Containerization: Docker Compose</li>
    <li>CI/CD: GitHub Actions</li>
    <li>Build tooling: Makefile</li>
  </ul>
</p>
<hr>
<h2>Setup Instructions</h2>
<pre>
You need Docker installed. Nothing else is required — Go, Node, and Postgres all run inside containers.

1. Clone the repository

bashgit clone https://github.com/anshikagupta17/MVC_Assignment.git
cd MVC_Assignment

2. Create your environment file

Copy the example file and fill in your own values:

bashcp .env.example .env

.env should contain:

DB_URI=postgres://username:password@localhost:5432/dbname?sslmode=disable
JWT_SECRET=your_secret_key_here

DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=assignment_mvc
SEED_DB=true

3. Start everything with Docker Compose

bashdocker compose up --build

This will:
Start a PostgreSQL container
Build and start the Go backend (automatically running all migrations, then seeding test data)
Build and start the Next.js frontend

To run it in the background instead:

bashdocker compose up --build -d

4. Access the app

Frontend: http://localhost:3000
Backend API: http://localhost:8080

5. Stopping everything

bashdocker compose down

To also wipe the database volume (full reset):

bashdocker compose down -v
</pre>



