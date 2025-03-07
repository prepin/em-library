from fastapi import FastAPI, HTTPException, status

app = FastAPI()


@app.get("/info")
def ok(song: str | None = None, group: str | None = None):
    if song is None or group is None:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Both 'name' and 'group' query parameters are required",
        )

    return {
        "releaseDate": "16.07.2006",
        "text": r"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
        "link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
    }


@app.get("/fail", status_code=status.HTTP_500_INTERNAL_SERVER_ERROR)
def fail():
    return {"detail": "everything failed"}
