
## Task 4

After playing around with profiling, I detected that storing the json in database, marshaling and unmarshaling it has some costs on throughput:

![Screenshot 1](photos/jsonops.png)  


Mostly dynamic allocations has their logic because of uncertainity of data to come. The only change ( a very little and banal one ) would be allocating the size of random domains read from json if known before processing. 

![Screenshot 1](photos/1mod.png)  

Before: 
![Screenshot 1](photos/1.png)  

After: 
![Screenshot 1](photos/1after.png)  


Well, the performance is quite affected as well by json datatype: 

![Screenshot 1](photos/perf.png)   


##### Storing the domains as json would easily give us the oppurtunity to decode structure that later could optimize the bussiness logic. Perhaps, we could cache the database responses to avoid the overhead.










