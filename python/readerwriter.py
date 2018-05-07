""" Read a line from a file, post it to service, read line, write it to file
"""
import asyncio
import click
import os
import requests  # just for the codes for now
import aiohttp
from urllib.parse import urljoin

queue_hostname = 'http://localhost:8080/'


@click.command()
@click.argument('inp', type=click.Path(exists=True))
@click.argument('out', type=click.Path())
def main(inp, out):
    loop = asyncio.get_event_loop()
    loop.set_exception_handler(None)
    path_name = os.path.basename(inp)
    queue_url = urljoin(queue_hostname, path_name)
    tasks = [
        asyncio.ensure_future(read_input(queue_url, inp)),
        asyncio.ensure_future(write_output(queue_url, out)),
    ]
    loop.run_until_complete(asyncio.gather(*tasks))
    loop.close()


async def read_input(queue_url, inp):
    async with aiohttp.ClientSession() as session:
        print('Using ', queue_url)
        r = await session.put(queue_url)
        r.raise_for_status()
        with open(inp, 'rb') as f:
            for l in f:
                async with session.post(queue_url, data=l) as r:
                    r.raise_for_status()
        r = await session.delete(queue_url)
        r.raise_for_status()


async def write_output(queue_url, out):
    with open(out, 'wb') as f:
        async with aiohttp.ClientSession() as session:
            while True:
                r = await get_line(session, queue_url)
                if r is None:
                    break
                if r == -1:
                    await asyncio.sleep(0.1)
                else:
                    f.write(r)


async def get_line(session, queue_url):
    async with session.get(queue_url) as r:
        if r.status == requests.codes.ok:
            return await r.read()
        elif r.status == requests.codes.GONE:
            return None
        elif r.status == requests.codes.no_content:
            return -1
        else:
            print(r)
            return None

if __name__ == '__main__':
    import logging
    logging.basicConfig(level=logging.DEBUG)
    main()
