from fastapi import FastAPI,UploadFile,status
from fastapi.requests import Request
from fastapi.middleware.cors import CORSMiddleware
from fastapi.staticfiles import StaticFiles
from fastapi.responses import JSONResponse
from pydantic_settings import BaseSettings
from fastapi_csrf_protect import CsrfProtect
from fastapi_csrf_protect.exceptions import CsrfProtectError
from fastapi import Depends
from dotenv import load_dotenv,find_dotenv
from PIL import Image
import os,clamd,hashlib,datetime,ffmpeg

cd = clamd.ClamdNetworkSocket(port=13310)

app = FastAPI()

load_dotenv(find_dotenv("cdn.env"),override=True)

if os.getenv("CSRF_SECRET") == None:
  print("Failed to load csrf protection!!")
  exit(1)

class CsrfSettings(BaseSettings):
  secret_key: str = os.getenv("CSRF_SECRET") or ""
  cookie_samesite: str = "lax"
  
@CsrfProtect.load_config
def get_csrf_config():
  return CsrfSettings()

app.add_middleware(
  CORSMiddleware,
  allow_origins=[
    "https://localhost:4000",
  ],
  allow_credentials=True,
  allow_methods=["GET","POST"],
  allow_headers=["*"],
)

@app.middleware("http")
async def authorization(request: Request, call_next):
  if request.headers.get("X-Token") == None or request.headers.get("X-Token") == os.getenv("API_KEY"):
    return JSONResponse({"status":"failed","result":"errors.unAuthorized"},status.HTTP_401_UNAUTHORIZED)
  return call_next(request)

app.mount("/files",StaticFiles(directory="files"),name="files")

@app.get("/upload")
def get_upload(csrf_protect: CsrfProtect = Depends()):
  csrf_token, signed_token = csrf_protect.generate_csrf_tokens()
  response = JSONResponse({"status":"success","token":csrf_token},200)
  csrf_protect.set_csrf_cookie(signed_token, response)
  return response

@app.post("/upload")
async def upload_file(request: Request, file: UploadFile, csrf_protect: CsrfProtect = Depends()):
  await csrf_protect.validate_csrf(request)
  if file.size/1024/1024 > 5 and file.content_type.startswith("image/"):
    response = JSONResponse({"status":"failed","result":"errors.largeFile"},status.HTTP_413_REQUEST_ENTITY_TOO_LARGE)
    csrf_protect.unset_csrf_cookie(response)
    return response
  content = await file.read()
  virusres = cd.instream(content)
  print(virusres["stream"][0])
  if virusres["stream"][0] == "FOUND":
    response = JSONResponse({"status":"failed","result":"errors.infectedFile"},status.HTTP_400_BAD_REQUEST)
    csrf_protect.unset_csrf_cookie(response)
    return response
  filename_without_ext = hashlib.md5(file.filename)+hashlib.md5(datetime.datetime.now().strftime("%d/%m/%Y,%H:%M:%S"))
  final_file = f"files/{filename_without_ext}.{file.filename.split(".")[-1]}"
  with open(final_file,"wb+") as nfile:
    nfile.write(file.read())
  if file.content_type.startswith("image/"):
    image = Image.open(final_file)
    data = list(image.getdata())
    nfile = Image.new(image.mode, image.size)
    nfile.putdata(data)
    nfile.save(f"files/{filename_without_ext}nmdicd.{file.filename.split(".")[-1]}",quality=70,lossless=True,optimize=True)
    final_file = f"files/{filename_without_ext}nmdicd.{file.filename.split(".")[-1]}"
  if file.content_type.startswith("video/"):
    (
      ffmpeg
      .input(final_file)
      .output(
        f"files/{filename_without_ext}nmdvcd.{file.filename.split(".")[-1]}",
        vf='scale=-1:720',
        vcodec='libx264',
        crf=23,
        preset='fast',
        acodec='aac',
        map_metadata=-1
      ).run()
    )
    final_file = f"files/{filename_without_ext}nmdvcd.{file.filename.split(".")[-1]}"
  
  response = JSONResponse()
  csrf_protect.unset_csrf_cookie(response)
  return response