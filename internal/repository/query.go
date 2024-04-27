package repository

const (
	// Users
	queryGetUserByParam = `
		SELECT 
			ID,
			Email,
			Password,
			Name,
			ProfilePicture,
			Active,
			Balance
		FROM Users
		WHERE %s = ?
		LIMIT 1;
	`
	qCreateUser = `
		INSERT INTO Users (
			Email, 
			Password, 
			Name 
		)
		VALUE (
			:Email, 
			:Password, 
			:Name
		);
	`
	qUpdateUserProfile = `
		UPDATE Users
		SET Name = :Name
		WHERE ID = :ID;
	`
	qUpdateUserField = `
		UPDATE Users
		SET %s = :Value
		WHERE ID = :ID;
	`
	qUpdateUserStatus = `
		UPDATE Users 
		SET Active = TRUE
		WHERE ID = ?;
	`
	qUpdateUserBalance = `
		UPDATE Users 
		SET Balance = Balance + :Value
		WHERE ID = :ID;
	`

	// UserVerifications
	qCreateUserVerification = `
		INSERT INTO UserVerifications (ID, UserID, Token)
		VALUE (:ID, :UserID, :Token);
	`
	qGetUserVerificationByIDAndToken = `
		SELECT ID, UserID, Token
		FROM UserVerifications
		WHERE ID = ? AND Token = ?
		LIMIT 1;
	`
	qUpdateUserVerificationStatus = `
		UPDATE UserVerifications
		SET Succeed = TRUE
		WHERE ID = ?;
	`

	// Notifications
	qFetchNotifications = `
		SELECT Text
		FROM Notifications
		WHERE ID = ?
		ORDER BY CreatedAt DESC
		LIMIT 5;
	`
	qCreateNotif = `
		INSERT INTO Notifications (UserID, Text)
		VALUE (:UserID, :Text);
	`

	// Enclosure
	qGetAllEnclosure = `
		SELECT Name, Coordinate
		FROM Enclosures;
	`
	qFetchDistanceToEnclosure = `
		SELECT ST_Distance_Sphere(Coordinate, POINT(:Longitude, :Latitude)) AS Distance 
		FROM Enclosures
		ORDER BY Distance ASC
		LIMIT 1;
	`
	// Exchanges
	qCreateExchange = `
		INSERT INTO Exchanges 
		(UserID, MerchantID, Amount, Date, Status)
		VALUE (:UserID, :MerchantID, :Amount, :Date, :Status);
	`
	qGetExchanges = `
		SELECT
			e.ID AS ID,
			e.UserID AS UserID,
			e.Amount AS Amount,
			e.Date AS Date,
			e.Status AS Status,
			m.ID AS MerchantID,
			m.Name AS MerchantName,
			m.Code AS MerchantCode
			FROM Exchanges e
			JOIN Merchants m ON e.MerchantID = m.ID
			WHERE %s  = ?
			ORDER BY e.CreatedAt DESC;
	`

	// Merchants
	qGetMerchantByParam = `
		SELECT
			ID,
			Name,
			Code
		FROM Merchants
		WHERE %s = ?
		LIMIT 1;
	`
)
