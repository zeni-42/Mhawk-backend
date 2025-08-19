package repository

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitTables(client *pgxpool.Pool) {
	if client == nil {
		log.Fatalf("DB client is nil can't proceed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exist bool
	err := client.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'users'
		);
	`).Scan(&exist)

	if err != nil {
		log.Fatalf("Failed to check if 'users' table exists: %v", err)
	}

	if !exist {
		userQuery :=
			`CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			fullname TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			refresh_token TEXT,
			is_new BOOLEAN DEFAULT TRUE,
			avatar TEXT DEFAULT 'https://res.cloudinary.com/dfbtssuwy/image/upload/v1735838884/ljziqvhelksqmytkffj9.jpg',
			is_pro BOOLEAN DEFAULT FALSE,
			is_organization BOOLEAN DEFAULT FALSE,
			organization_id UUID,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`

		apiQuery := `
			CREATE TABLE apikeys (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
				key_name TEXT UNIQUE NOT NULL, 
				api_key TEXT UNIQUE NOT NULL,
				description TEXT NOT NULL,
				is_active BOOLEAN DEFAULT TRUE,
				environment TEXT NOT NULL, 
				token INTEGER DEFAULT 50,
				used_token INTEGER DEFAULT 0,
				expired_date TIMESTAMPTZ,
				refresh_date TIMESTAMPTZ,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW()
			);`

		emailQuery := `
			CREATE TABLE emails (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID REFERENCES users(id) ON DELETE CASCADE,
				api_id	UUID REFERENCES apikeys(id) ON DELETE CASCADE,
				"to" TEXT NOT NULL,
				subject TEXT NOT NULL,
				body TEXT NOT NULL,
				html TEXT NOT NULL,
				is_html BOOLEAN DEFAULT FALSE,
				status TEXT DEFAULT 'Pending',
				is_template BOOLEAN DEFAULT FALSE,
				template_id UUID,
				is_bulk BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW()
			);`

		organizationQuery := `
			CREATE TABLE organization (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				founder UUID REFERENCES users(id) ON DELETE CASCADE,
				name TEXT NOT NULL,
				description TEXT NOT NULL,
				domain TEXT NOT NULL,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW()
			);`

		organizationMembersQuery := `
			CREATE TABLE organization_members (
				organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
				user_id UUID REFERENCES users(id) ON DELETE CASCADE,
				PRIMARY KEY (organization_id, user_id)
			);`

		organizationApiKeyQuery := `
			CREATE TABLE organization_api_keys (
				organization_id UUID REFERENCES organization(id) ON DELETE CASCADE,
				api_id UUID REFERENCES apikeys(id) ON DELETE CASCADE,
				PRIMARY KEY (organization_id, api_id)
			);`

		templateQuery := `
			CREATE TABLE template (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				name 	TEXT,
				subject TEXT,
				body	TEXT,
				author_id UUID,
				organization_id UUID,
				is_restricted BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW()
			);`

		logQuery := `
			CREATE TABLE logs (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				email_id UUID REFERENCES emails(id) ON DELETE CASCADE,
				status TEXT DEFAULT 'queued',
				created_at TIMESTAMPTZ DEFAULT NOW(),
				updated_at TIMESTAMPTZ DEFAULT NOW()
			);`

		notificationsQuery := `
			CREATE TABLE notifications (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
				reciever_id UUID REFERENCES users(id) ON DELETE CASCADE,
				is_system BOOLEAN DEFAULT FALSE,
				title TEXT NOT NULL,
				message TEXT NOT NULL,
				data JSONB NOT NULL,
				is_read BOOLEAN DEFAULT FALSE,
				read_at TIMESTAMPTZ,
				created_at TIMESTAMPTZ DEFAULT NOW()
			);`

		if _, err := client.Exec(ctx, userQuery); err != nil {
			log.Fatalf("Failed to create 'users' table ==> %v", err)
		}

		if _, err := client.Exec(ctx, apiQuery); err != nil {
			log.Fatalf("Failed to create 'api_key' table ==> %v", err)
		}

		if _, err := client.Exec(ctx, organizationQuery); err != nil {
			log.Fatalf("Failed to create 'organization' table ==> %v", err)
		}

		if _, err := client.Exec(ctx, organizationMembersQuery); err != nil {
			log.Fatalf("Failed to create 'organization_members' table ==> %v", err)
		}

		if _, err := client.Exec(ctx, organizationApiKeyQuery); err != nil {
			log.Fatalf("Failed to create 'organization_api_keys' table ==> %v", err)
		}

		if _, err := client.Exec(ctx, emailQuery); err != nil {
			log.Fatalf("Failed to create 'emails' table ==> %v", err)
		}
		if _, err := client.Exec(ctx, logQuery); err != nil {
			log.Fatalf("Failed to create 'logs' table ==> %v", err)
		}

		if _, err := client.Exec(ctx, templateQuery); err != nil {
			log.Fatalf("Failed to create 'template' table ==> %v", err)
		}

		if _, err := client.Exec(ctx, notificationsQuery); err != nil {
			log.Fatalf("Failed to create 'notifications' table ==> %v", err)
		}

	}
}
