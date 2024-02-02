package main

import (
    "fmt"
    "github.com/minio/minio-go/v7"
    "log"
    "os"

    "context"
    "github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
    endpoint := "114.215.190.173:9000"
    //endpoint := "play.min.io"
    accessKeyID := "minioadmin"
    secretAccessKey := "xx123456"

    // Initialize minio client object.
    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
    })
    if err != nil {
        log.Fatalln(err)
    }

    //log.Printf("%#v\n", minioClient) // minioClient is now setup

    buckets, err := minioClient.ListBuckets(context.Background())
    if err != nil {
        fmt.Println(err)
        return
    }
    for _, bucket := range buckets {
        fmt.Println(bucket)
    }
    //upload(minioClient)
    get(minioClient)
}

func upload(minioClient *minio.Client) {
    file, err := os.Open("14310208948_1679453569.amr")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()

    fileStat, err := file.Stat()
    if err != nil {
        fmt.Println(err)
        return
    }

    uploadInfo, err := minioClient.PutObject(context.Background(), "test", "myobject", file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("Successfully uploaded bytes: ", uploadInfo)
}

func get(minioClient *minio.Client) {
    object, err := minioClient.GetObject(context.Background(), "test", "myobject", minio.GetObjectOptions{})
    if err != nil {
        fmt.Println(err)
        return
    }
    defer object.Close()

    databuf := make([]byte, 40960)
    object.Read(databuf)
    fmt.Println(string(databuf))
}
