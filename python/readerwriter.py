""" Read a line from a file, post it to service, read line, write it to file
"""
import asyncio
import click
import requests

queue_hostname = 'http://localhost:8080/'


@click.command()
@click.argument('inp')
@click.argument('out')
def main(inp, out):
    loop = asyncio.get_event_loop()
    tasks = [
        asyncio.ensure_future(read_input(inp)),
        asyncio.ensure_future(write_output(out)),
    ]
    loop.run_until_complete(asyncio.wait(tasks))
    loop.close()



async def read_input(inp):
    with open(inp) as f:
        for l in f:
            r = requests.post(queue_hostname, data=l)
            r.raise_for_status()


async def write_output(out):
    with open(out, 'wb') as f:
        while True:
            r = requests.get(queue_hostname)
            if r.status_code == 200:
                f.write(r.content)
            else:
                print(r)
                break

if __name__ == '__main__':
    import logging
    logging.basicConfig(level=logging.DEBUG)
    main()
