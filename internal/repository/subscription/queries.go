package subscription

const (
	subscriptionColumns = `
		id,
		service_name,
		price,
		user_id,
		start_date,
		end_date,
		created_at,
		updated_at`

	createSubscription = `
		INSERT INTO subscriptions (
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		RETURNING ` + subscriptionColumns

	getSubscription = `
		SELECT ` + subscriptionColumns + `
		FROM subscriptions
		WHERE id = $1`

	updateSubscription = `
		UPDATE subscriptions
		SET
			service_name = $2,
			price = $3,
			user_id = $4,
			start_date = $5,
			end_date = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING ` + subscriptionColumns

	deleteSubscription = `
		DELETE FROM subscriptions
		WHERE id = $1`

	listSubscriptions = `
		SELECT ` + subscriptionColumns + `
		FROM subscriptions
		WHERE ($1::uuid IS NULL OR user_id = $1)
		  AND ($2::text IS NULL OR service_name ILIKE '%' || $2 || '%')
		ORDER BY created_at DESC
		LIMIT $3
		OFFSET $4`
)
