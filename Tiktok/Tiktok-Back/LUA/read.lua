-- read.lua
-- 95% 请求落在 1000 个热门视频，制造热 key
local ids = {}
for i = 1, 1000 do
    ids[i] = "video_hot_" .. i
end

request = function()
    local video_id = ids[math.random(1, 1000)]
    local path = "/api/videos/get-video-likes?video_id=" .. video_id
    return wrk.format("GET", path)
end