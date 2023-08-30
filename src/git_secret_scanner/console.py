from __future__ import annotations
from typing import cast, Self, Any
from types import TracebackType

from rich import console, progress
import sys


stdout = console.Console()
stderr = console.Console(stderr=True)


def print(message: str) -> None:  # noqa: A001
    stdout.print(message)


def print_error(message: str, error: Exception | None = None) -> None:
    content = message
    if error:
        content = f'{message}: {error}'
        if hasattr(error, '__notes__') and len(error.__notes__) > 0:
            for note in error.__notes__:
                content += f'\n\n{note}'
    stderr.print(f'[red]{content}[/red]')


def exit_with_error(message: str, error: Exception | None = None) -> None:
    print_error(message, error)
    sys.exit(1)


class Progress(progress.Progress):
    def __init__(self: Self,
        *args: Any,  # noqa: ANN401
        description: str = '',
        transient: bool = False,
    ) -> None:
        super().__init__(*args, console=stdout, transient=transient)
        self.description = description

    def __exit__(self: Self,
        exc_type: type[BaseException] | None,
        exc_val: BaseException | None,
        exc_tb: TracebackType | None,
    ) -> bool:
        super().__exit__(exc_type, exc_val, exc_tb)
        if exc_type is None and exc_val is None:
            print(f'[green]✔[/green] {self.description}')
            return True
        print(f'[red]⨯[/red] {self.description}')  # noqa: RUF001
        return False


class ProgressSpinner(Progress):
    def __init__(self: Self, description: str) -> None:
        super().__init__(
            progress.SpinnerColumn(),
            progress.TextColumn('[progress.description]{task.description}'),
            transient=True,
            description=description,
        )
        self.task = super().add_task(description=description, total=None)

    # useful for typing
    def __enter__(self: Self) -> Self:
        return cast(ProgressSpinner, super().__enter__())


class ProgressBar(Progress):
    def __init__(self: Self, description: str, num_steps: int) -> None:
        super().__init__(transient=True, description=description)
        self.task = super().add_task(description, total=num_steps)

    def update(self: Self, advance: float) -> None:
        super().update(self.task, advance=advance)

    # useful for typing
    def __enter__(self: Self) -> Self:
        return cast(ProgressBar, super().__enter__())
