package models

//ServicePing defines a row in service_ping Table
type ServicePing struct {
	ID           int    `pg:"id,pk"`
	Name         string `pg:"name"`
	Timestamp    int64  `pg:"ping_timestamp"`
	IsAvailable  bool   `pg:"is_available"`
	ResponseTime int    `pg:"response_time"`
}

//A better solution could have been to define a sperate table of service (id, name(unique), url(unique))
//And Another Table for Ping records which takes a service_id as foreign key
//But went for a simpler solution to save time as this is just to make sure I know
//How things work
