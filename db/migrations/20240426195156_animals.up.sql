CREATE TABLE Animals (
    ID BIGINT UNSIGNED AUTO_INCREMENT,
    
    Picture TEXT NOT NULL,
    Name VARCHAR(255) NOT NULL,
    Latin VARCHAR(255) NOT NULL,
    Origin VARCHAR(255) NOT NULL,
    Characteristic VARCHAR(255) NOT NULL,
    Diet VARCHAR(255) NOT NULL,
    Lifespan VARCHAR(255) NOT NULL,
    EnclosureCoordinate POINT NOT NULL,
    
    CreatedAt DATETIME NOT NULL DEFAULT NOW(),
    UpdatedAt DATETIME NOT NULL DEFAULT NOW(),
    
    PRIMARY KEY (ID)
) ENGINE = INNODB DEFAULT CHARSET = UTF8;

INSERT INTO Animals(Picture, Name, Latin, Origin, Characteristic, Diet, Lifespan, EnclosureCoordinate) VALUES
("https://aoyfqtbcektucqlryxvj.supabase.co/storage/v1/object/public/bucket/harimau.jpg?t=2024-04-27T12%3A44%3A04.744Z", "Harimau Sumatera", "Panthera Tigris", "Sumatera, Indonesia", "Harimau Sumatera adalah salah satu subspesies harimau yang hanya ditemukan di habitat asli pulau Sumatera, yang merupakan bagian dari Indonesia.", "Karnivora", "20 Tahun", POINT(112.73602656088153, -7.296691437475779)),
("https://aoyfqtbcektucqlryxvj.supabase.co/storage/v1/object/public/bucket/gajah.jpg?t=2024-04-27T12%3A45%3A04.476Z", "Gajah Sumatera", "Elephas maximus sumatrensis", "Sumatera, Indonesia", "Gajah Sumatera adalah subspesies gajah yang terancam punah dan endemik di pulau Sumatera, Indonesia.  Gajah Sumatera mendiami hutan hujan tropis dan lahan gambut di Sumatera.", "Herbivora", "35 Tahun", POINT(112.73633233270644, -7.297286721241508));
