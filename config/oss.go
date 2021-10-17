package config

// basic oss config

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssConfig struct {
	Endpoint  string `json:"endpoint" yaml:"endpoint" toml:"endpoint" xml:"endpoint"`
	AccessID  string `json:"access_id" yaml:"access_id" toml:"access_id" xml:"access_id"`
	AccessKey string `json:"access_key" yaml:"access_key" toml:"access_key" xml:"access_key"`
	Bucket    string `json:"bucket" yaml:"bucket" toml:"bucket" xml:"bucket"`
}

func (c *OssConfig) Client() (*oss.Client, error) {
	return oss.New(c.Endpoint, c.AccessID, c.AccessKey)
}

func (c *OssConfig) ClientG() (*oss.Client, error) {
	client, err := oss.New(c.Endpoint, c.AccessID, c.AccessKey)
	if err != nil {
		return nil, err
	}
	buk, err := client.Bucket(c.Bucket)
	if err != nil {
		return nil, err
	}
	Global.FileStore = NewOSSStore(buk)
	return client, nil
}

func UploadDir(bucket *oss.Bucket, localdir, ossBaseDir string) error {
	info, err := ioutil.ReadDir(localdir)
	if err != nil {
		return err
	}
	for _, v := range info {
		if v.IsDir() {
			err := UploadDir(bucket, path.Join(localdir, v.Name()), ossBaseDir)
			if err != nil {
				return err
			}
		} else {
			ossP := path.Join(ossBaseDir, localdir, v.Name())
			t, err := bucket.GetObjectTagging(ossP)
			if err == nil {
				var ok bool
				for _, tag := range t.Tags {
					if tag.Key == "md5" {
						ok = true
						break
					}
				}
				if ok {
					continue
				}
			}
			f, err := os.Open(path.Join(localdir, v.Name()))
			if err != nil {
				return err
			}
			h := md5.New()
			r := io.TeeReader(f, h)
			err = bucket.PutObject(ossP, r)
			if err != nil {
				f.Close()
				return err
			}
			f.Close()

			err = bucket.PutObjectTagging(ossP, oss.Tagging{
				Tags: []oss.Tag{
					{
						Key:   "md5",
						Value: fmt.Sprintf("%x", h.Sum(nil)),
					},
				},
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
