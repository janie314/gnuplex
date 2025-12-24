-- Remove the item from the playlist when it finishes playing
mp.register_event("end-file", function()
    mp.command("playlist-remove current")
end)
