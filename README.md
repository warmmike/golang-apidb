# golang-apidb
Golang Project API DB

![alt text](https://github.com/warmmike/golang-apidb/blob/main/golang_project_api_db.png?raw=true)

## Requirements
- API, Golang, Kubernetes workload; data queried by year, movie name, cast member or genre
- 24/7 zero downtime
- 5ms response time
- Diagram
- [Data format](https://github.com/prust/wikipedia-movie-data/tree/master)

## Examples
By year
```
mike@Michaels-MacBook-Pro golang-apidb % curl -X GET 'http://localhost:8081/movies/?year=1990'
[{"title":"Arachnophobia","year":1990,"cast":["Jeff Daniels","John Goodman","Harley Jane Kozak","Julian Sands"],"genres":["Thriller","Comedy","Horror"],"href":"Arachnophobia_(film)","extract":"Arachnophobia is a 1990 American horror comedy film directed by Frank Marshall in his directorial debut from a screenplay by Don Jakoby and Wesley Strick. Starring Jeff Daniels and John Goodman, the film follows a small California town that becomes invaded by an aggressive and dangerous spider species. Its title refers to the fear of spiders.","thumbnail":"https://upload.wikimedia.org/wikipedia/en/a/a0/Arachnophobia_%28film%29_POSTER.jpg","thumbnail_width":259,"thumbnail_height":383}]
```

By movie name
```
mike@Michaels-MacBook-Pro golang-apidb % curl -X GET 'http://localhost:8081/movies/?title=Popeye'
[{"title":"Popeye","year":1980,"cast":["Robin Williams","Shelley Duvall","Paul Dooley"],"genres":["Musical","Comedy"],"href":"Popeye_(film)","extract":"Popeye is a 1980 American musical comedy film directed by Robert Altman and produced by Paramount Pictures and Walt Disney Productions. It is based on E. C. Segar's Popeye comics character. The script was written by Jules Feiffer, and stars Robin Williams as Popeye the Sailor Man and Shelley Duvall as Olive Oyl. Its story follows Popeye's adventures as he arrives in the town of Sweethaven.","thumbnail":"https://upload.wikimedia.org/wikipedia/en/8/88/Popeyemovieposter.jpg","thumbnail_width":258,"thumbnail_height":391}]
```

By cast member
```
mike@Michaels-MacBook-Pro golang-apidb % curl -X GET 'http://localhost:8081/movies/?cast=Chris%20Evans'       
[{"title":"The Avengers","year":2012,"cast":["Robert Downey Jr.","Chris Evans","Chris Hemsworth","Mark Ruffalo","Jeremy Renner","Scarlett Johansson","Tom Hiddleston","Samuel L. Jackson","Stellan Skarsgård","Cobie Smulders","Clark Gregg","Gwyneth Paltrow","Maximiliano Hernández","Paul Bettany","Alexis Denisof","Damion Poitier","Powers Boothe","Jenny Agutter","Stan Lee","Harry Dean Stanton","Jerzy Skolimowski","Warren Kole","Enver Gjokaj"],"genres":["Superhero"],"href":"The_Avengers_(2012_film)","extract":"Marvel's The Avengers, or simply The Avengers, is a 2012 American superhero film based on the Marvel Comics superhero team of the same name. Produced by Marvel Studios and distributed by Walt Disney Studios Motion Pictures, it is the sixth film in the Marvel Cinematic Universe (MCU). Written and directed by Joss Whedon, the film features an ensemble cast including Robert Downey Jr., Chris Evans, Mark Ruffalo, Chris Hemsworth, Scarlett Johansson, and Jeremy Renner as the Avengers, alongside Tom Hiddleston, Stellan Skarsgård, and Samuel L. Jackson. In the film, Nick Fury and the spy agency S.H.I.E.L.D. recruit Tony Stark, Steve Rogers, Bruce Banner, Thor, Natasha Romanoff, and Clint Barton to form a team capable of stopping Thor's brother Loki from subjugating Earth.","thumbnail":"https://upload.wikimedia.org/wikipedia/en/8/8a/The_Avengers_%282012_film%29_poster.jpg","thumbnail_width":220,"thumbnail_height":326},{"title":"Avengers: Age of Ultron","year":2015,"cast":["Robert Downey Jr.","Chris Hemsworth","Mark Ruffalo","Chris Evans","Scarlett Johansson","Jeremy Renner","Don Cheadle","Aaron Taylor-Johnson","Elizabeth Olsen","Paul Bettany","Cobie Smulders","Anthony Mackie","Hayley Atwell","Idris Elba","Stellan Skarsgård","James Spader","Samuel L. Jackson"],"genres":["Superhero"],"href":"Avengers:_Age_of_Ultron","extract":"Avengers: Age of Ultron is a 2015 American superhero film based on the Marvel Comics superhero team the Avengers. Produced by Marvel Studios and distributed by Walt Disney Studios Motion Pictures, it is the sequel to The Avengers (2012) and the 11th film in the Marvel Cinematic Universe (MCU). Written and directed by Joss Whedon, the film features an ensemble cast including Robert Downey Jr., Chris Hemsworth, Mark Ruffalo, Chris Evans, Scarlett Johansson, Jeremy Renner, Don Cheadle, Aaron Taylor-Johnson, Elizabeth Olsen, Paul Bettany, Cobie Smulders, Anthony Mackie, Hayley Atwell, Idris Elba, Linda Cardellini, Stellan Skarsgård, James Spader, and Samuel L. Jackson. In the film, the Avengers fight Ultron (Spader)—an artificial intelligence created by Tony Stark (Downey) and Bruce Banner (Ruffalo) who plans to bring about world peace by causing human extinction.","thumbnail":"https://upload.wikimedia.org/wikipedia/en/f/ff/Avengers_Age_of_Ultron_poster.jpg","thumbnail_width":220,"thumbnail_height":326}]
```

Be genre
```
mike@Michaels-MacBook-Pro golang-apidb % curl -X GET 'http://localhost:8081/movies/?genre=Comedy'
[{"title":"Popeye","year":1980,"cast":["Robin Williams","Shelley Duvall","Paul Dooley"],"genres":["Musical","Comedy"],"href":"Popeye_(film)","extract":"Popeye is a 1980 American musical comedy film directed by Robert Altman and produced by Paramount Pictures and Walt Disney Productions. It is based on E. C. Segar's Popeye comics character. The script was written by Jules Feiffer, and stars Robin Williams as Popeye the Sailor Man and Shelley Duvall as Olive Oyl. Its story follows Popeye's adventures as he arrives in the town of Sweethaven.","thumbnail":"https://upload.wikimedia.org/wikipedia/en/8/88/Popeyemovieposter.jpg","thumbnail_width":258,"thumbnail_height":391},{"title":"Arachnophobia","year":1990,"cast":["Jeff Daniels","John Goodman","Harley Jane Kozak","Julian Sands"],"genres":["Thriller","Comedy","Horror"],"href":"Arachnophobia_(film)","extract":"Arachnophobia is a 1990 American horror comedy film directed by Frank Marshall in his directorial debut from a screenplay by Don Jakoby and Wesley Strick. Starring Jeff Daniels and John Goodman, the film follows a small California town that becomes invaded by an aggressive and dangerous spider species. Its title refers to the fear of spiders.","thumbnail":"https://upload.wikimedia.org/wikipedia/en/a/a0/Arachnophobia_%28film%29_POSTER.jpg","thumbnail_width":259,"thumbnail_height":383},{"title":"Herbie Goes Bananas","year":1980,"cast":["Cloris Leachman","Harvey Korman"],"genres":["Comedy","Adventure"],"href":"Herbie_Goes_Bananas","extract":"Herbie Goes Bananas is a 1980 American adventure comedy film directed by Vincent McEveety and written by Don Tait. The film is the fourth installment in the Herbie film series and the sequel to Herbie Goes to Monte Carlo (1977).","thumbnail":"https://upload.wikimedia.org/wikipedia/en/e/e1/Herbie_goes_bananas_poster.jpg","thumbnail_width":255,"thumbnail_height":390}]
```