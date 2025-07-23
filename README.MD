# Description

Pipeline visualizer

Backend has mock runner implementation, which executes pipelines

Frontend has visualization of pipeline with controls

https://github.com/user-attachments/assets/3f510751-8f88-4ac4-a609-7ef995292583

# Fast start
 
1) Clone
2) docker compose up --build
3) Login. Login **admin**, password **admin**
5) Open **localhost:8080**
6) Check functionality
7) Another user. Login **tester** password **tester**. Has read permissions, but not execution 

# Features
- Working mock runner
- Simple and minimalistic interface
- Pipeline statistics, updating in real-time
- Can export statistics to excel
- Jobs sorting, minimizing cross-crossings of depenency connections
- Scalable visualization, supporting many jobs
- Database is used for information storing
- Job progress updates in real time
- Very flexible controls of jobs. You can control each job separately, not necessarily entire pipeline
- Supports tablets and phones, adaptive inferface and controls
- Read and execute permissions for pipelines

# Tech stack
- NextJS - frontend, javascript, no libraries (except for date-fns)
- Go - backend, gorilla and gorm
- Postgres - database
