import asyncio
import aiofiles


async def rea():
    while True:
        file = await aiofiles.open("/tmp/fif", "r")
        while out:=await file.readline():
            out=out.rstrip()
            print(out)
           
async def timer():
    while True:
        print("aaa")
        await asyncio.sleep(5)


async def main():
    await asyncio.gather(rea(), timer())
       
asyncio.run(main())
