import "time"

#Metadata: {
	title:    string
	status:   #Status
	owner:    #Owner
	reviewed: #Date
	tags: [...string]
}

#Status: "archived" | "draft" | "published"

#Date: time.Format(time.RFC3339Date)

#Owner: string
