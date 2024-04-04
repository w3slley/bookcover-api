package handlers

type GoogleBook struct{
  Kind string `json:"kind"`
  TotalItems int `json:"totalItems"` 
  Items []BookItem `json:"items"`
}

type BookItem struct {
  VolumeInfo BookInfo `json:"volumeInfo"`
}

type BookInfo struct {
  ImageLinks ImageLinks `json:"imageLinks"`
}

type ImageLinks struct {
  SmallThumbnail string `json:"smallThumbnail"`
  Thumbnail string `json:"thumbnail"`
}
