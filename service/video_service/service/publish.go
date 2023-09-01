package video_service

const publicPathPre = "../../public/"
const avatarPathPre = publicPathPre + "Avatar/"
const backgroundImagePathPre = publicPathPre + "BackgroudImage/"
const coverPathPre = publicPathPre + "Cover/"
const videoPathPre = publicPathPre + "Video/"

// UploadVideo 上传视频 放到 /public/Video 里
func uploadVideo(data []byte, videoTitle string) (string, error) {
	//	TODO
	return "", nil
}

// UploadCover 截取并上次封面 放到 /public/Cover 里
func uploadCover(playUrl string, coverTitle string) error {
	//	TODO
	return nil
}

// VideoPublish 上传视频
func VideoPublish(data []byte, videoTitle string, coverTitle string) error {
	playUrl, err := uploadVideo(data, videoTitle)
	if err != nil {
		return err
	}
	err = uploadCover(playUrl, coverTitle)
	if err != nil {
		return err
	}
	return nil
}
