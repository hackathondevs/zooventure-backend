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
			Active
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
		WHERE ID = ?;
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
)
