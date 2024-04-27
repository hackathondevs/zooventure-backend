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
			Balance,
			Admin
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

	// Campaign
	qGetAllCampaignByUserID = `
		SELECT 
		Campaigns.ID, 
		Campaigns.Picture, 
		Campaigns.Title, 
		Campaigns.Reward, 
		CASE WHEN CampaignSubmissions.Submission IS NOT NULL THEN TRUE ELSE FALSE AS Submitted
		FROM Campaigns
		LEFT JOIN CampaignSubmissions
		ON CampaignSubmissions.CampaignID = Campaigns.ID
		WHERE CampaignSubmissions.UserID = ?
		ORDER BY CreatedAt DESC;
	`
	qGetCampaignWithSubmission = `
		SELECT 
		Campaigns.Picture, 
		Campaigns.Title, 
		Campaigns.Description, 
		Campaigns.Reward,
		CASE WHEN CampaignSubmissions.Submission IS NOT NULL THEN TRUE ELSE FALSE AS Submitted
		CampaignSubmissions.Submission
		FROM Campaigns
		LEFT JOIN CampaignSubmissions
		ON CampaignSubmissions.CampaignID = Campaigns.ID
		WHERE Campaigns.ID = ? AND CampaignSubmissions.UserID = ?
		LIMIT 1;
	`
	qCreateCampaign = `
		INSERT INTO Campaigns 
		(Picture, Title, Description, Reward)
		VALUE 
		(:Picture, :Title, :Description, :Reward);
	`
	qUpdateCampaign = `
		UPDATE 
			Campaigns
		SET 
			Picture = :Picture, 
			Title = :Title, 
			Description = :Description, 
			Reward = :Reward
		WHERE ID = :ID;
	`
	qDeleteCampaign = `
		DELETE FROM Campaigns
		WHERE ID = ?;
	`
	qGetCampaignByID = `
		SELECT
			ID,
			Picture,
			Title,
			Description,
			Reward
		FROM Campaigns
		WHERE ID = ?
		LIMIT 1;
	`
	qGetAllCampaign = `
		SELECT
			ID,
			Picture,
			Title,
			Description,
			Reward
		FROM Campaigns
	`
	qCreateSubmission = `
		INSERT INTO CampaignSubmissions
		(CampaignID, UserID, Submission)
		VALUE
		(:CampaignID, :UserID, :Submission);`

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
	`

	// Reports
	qCreateReport = `
		INSERT INTO Reports
		(UserID, Picture, Description, Location)
		VALUE
		(:UserID, :Picture, :Description, :Location);
	`
	qGetReports = `
		SELECT
			ID,
			UserID,
			Picture,
			Description,
			Location,
			CreatedAt,
			Action
		FROM Reports
		ORDER BY CreatedAt DESC;
	`
	qUpdateReport = `
		UPDATE Reports
		SET
			Action = :Action
		WHERE ID = :ID;
	`
	qGetReportByID = `
		SELECT
			ID,
			UserID,
			Picture,
			Description,
			Location,
			CreatedAt,
			Action
		FROM Reports
		WHERE ID = ?
	`
	// Animals
	qGetAllAnimals = `
	    SELECT,
			ID,
			Name,
			Latin,
			Origin,
			Characteristics,
			Diet,
			Lifespan,
			EnclosureCoordinate
		FROM Animals;
	`
	qFetchTopRelated = `
		SELECT
			ID,
			Picture,
			Name,
			Latin,
			Origin,
			Characteristic,
			Diet,
			Lifespan,
			EnclosureCoordinate,
			ST_Distance_Sphere(EnclosureCoordinate, POINT(:Longitude, :Latitude)) AS Distance
		FROM Animals
		WHERE Name LIKE :Name
		ORDER BY Distance ASC
		LIMIT 1;
	`
)
