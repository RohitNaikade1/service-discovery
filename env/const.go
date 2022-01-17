package env

var (
	DB_URL                  = "localhost:27017"
	DB_USERNAME             = "sdadmin"
	DB_PASSWORD             = "servicediscoverydev"
	MONGODB_DATABASE        = "service-discovery"
	ADMIN_USERNAME          = "sd_admin"
	ADMIN_PASSWORD          = "servicediscovery"
	ADMIN_EMAIL             = "sdadmin@gmail.com"
	PORT                    = GetEnvironmentVariable("PORT")
	USER_COLLECTION         = GetEnvironmentVariable("USER_COLLECTION")
	CREDENTIAL_COLLECTION   = GetEnvironmentVariable("CREDENTIAL_COLLECTION")
	REGISTRATION_COLLECTION = GetEnvironmentVariable("REGISTRATION_COLLECTION")
)
