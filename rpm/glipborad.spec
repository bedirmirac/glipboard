Name:           glipboard
Version:        2.1.1
Release:        1%{?dist}
Summary:        A terminal-based clipboard manager

License:        MIT
URL:            https://github.com/bedirmirac/glipboard
Source0:        %{url}/archive/v%{version}.tar.gz

BuildRequires:  golang >= 1.20

%description
Glipboard is an open-source clipboard manager written in Go using Bubble Tea and SQLite. 
It features event-driven copying and provides a terminal user interface (TUI) for managing 
clipboard history.

%prep
%autosetup -n %{name}-%{version}

%build
export CGO_ENABLED=0
export GOTOOLCHAIN=auto
go build -v -trimpath -o %{name} .

%install
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}

%changelog
* Tue Aug 18 2026 Miraç Bedir <eposta@adresin.com> - 2.1.1-1
- Update to version 2.1.1
- Build as a statically linked pure Go binary (CGO_ENABLED=0)
