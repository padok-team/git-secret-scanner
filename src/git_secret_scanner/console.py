from __future__ import annotations
import typing

from rich import console, progress


stdout = console.Console()
stderr = console.Console(stderr=True)


def print(message: str):
    stdout.print(message)


def print_error(message: str, error: Exception | None = None):
    content = message
    if error:
        content = f'{message}: {error}'
        if len(error.__notes__) > 0:
            for note in error.__notes__:
                content += f'\n\n{note}'
    stderr.print(f'[red]{content}[/red]')


def exit_with_error(message: str, error: Exception | None = None):
    print_error(message, error)
    exit(1)


class Progress(progress.Progress):
    def __init__(self, *args, description='', **kwargs):
        super().__init__(*args, **kwargs)
        self.description = description

    def stop(self):
        super().stop()
        print(f'[green]✔[/green] {self.description}')

    def error(self):
        super().stop()
        print(f'[red]⨯[/red] {self.description}')


class ProgressSpinner(Progress):
    def __init__(self, description: str):
        super().__init__(
            progress.SpinnerColumn(),
            progress.TextColumn('[progress.description]{task.description}'),
            transient=True,
            description=description,
        )
        self.task = super().add_task(description=description, total=None)

    # useful for typing
    def __enter__(self) -> ProgressSpinner:
        return typing.cast(ProgressSpinner, super().__enter__())


class ProgressBar(Progress):
    def __init__(self, description: str, num_steps: int):
        super().__init__(transient=True, description=description)
        self.task = super().add_task(description, total=num_steps)
    
    def update(self, advance: float):
        super().update(self.task, advance=advance)

    # useful for typing
    def __enter__(self) -> ProgressBar:
        return typing.cast(ProgressBar, super().__enter__())
