""" Read a line from a file, post it to service, read line, write it to file
"""
import click


@click.command()
def main():
    pass


if __name__ == '__main__':
    import logging
    logging.basicConfig(level=logging.DEBUG)
    main()
