import secrets
from typing import Annotated

from fastapi import Depends, FastAPI, HTTPException, status
from fastapi.responses import RedirectResponse
from pydantic import BaseModel, Field


class InMemoryStore:
	__slots__ = ("urls",)

	def __init__(self) -> None:
		self.urls: dict[str, str] = {}

	def set(self, code: str, url: str) -> None:
		self.urls[code] = url

	def get(self, code: str) -> str | None:
		return self.urls.get(code)


class ShortenRequest(BaseModel):
	url: Annotated[str, Field(min_length=1)]


class ShortenResponse(BaseModel):
	message: str
	short_url: str


app = FastAPI()
store = InMemoryStore()
base_url = "http://localhost:8080"


def get_store() -> InMemoryStore:
	return store


StoreDependency = Annotated[
	InMemoryStore,
	Depends(get_store),
]


def generate_short_code() -> str:
	return secrets.token_hex(4)


@app.get("/")
async def root() -> dict[str, str]:
	return {"message": "Hello world"}


@app.post(
	"/url",
	response_model=ShortenResponse,
	status_code=status.HTTP_201_CREATED,
)
async def shorten_url(
	request: ShortenRequest,
	store: StoreDependency,
) -> ShortenResponse:
	code = generate_short_code()
	store.set(code, request.url)

	return ShortenResponse(
		message="URL enregistrée",
		short_url=f"{base_url}/{code}",
	)


@app.get("/{code}")
async def redirect_to_url(
	code: str,
	store: StoreDependency,
) -> RedirectResponse:
	target_url = store.get(code)

	if target_url is None:
		raise HTTPException(status_code=status.HTTP_404_NOT_FOUND)

	return RedirectResponse(
		url=target_url,
		status_code=status.HTTP_307_TEMPORARY_REDIRECT,
	)