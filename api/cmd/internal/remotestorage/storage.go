package remotestorage

//Define the interfaces a remote storage implementation must satisfy

type Uploader interface {
	Upload(bucketName, fileName string) error
}

type Retriever interface {
	Retrieve(buckeName, filename string) error
}

type UploaderRetriever interface {
	Uploader
	Retriever
}
