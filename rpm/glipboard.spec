Name:           glipboard
Version:        2.1.1
Release:        3%{?dist}
Summary:        A terminal-based clipboard manager

License:        MIT
URL:            https://github.com/bedirmirac/glipboard
Source0:        %{url}/archive/v%{version}.tar.gz
Source1:        glipboard.desktop
Source2:        glipboard.service
Source3:        glipboard.png

BuildRequires:  golang >= 1.20
%{?systemd_requires}
BuildRequires:  systemd-rpm-macros

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

install -Dpm 0644 %{SOURCE1} %{buildroot}%{_datadir}/applications/%{name}.desktop

install -Dpm 0644 %{SOURCE2} %{buildroot}%{_userunitdir}/%{name}.service

install -Dpm 0644 %{SOURCE3} %{buildroot}%{_datadir}/icons/hicolor/256x256/apps/%{name}.png

%post
%systemd_user_post %{name}.service

%preun
%systemd_user_preun %{name}.service

%postun
%systemd_user_postun %{name}.service

%files
%license LICENSE
%doc README.md
%{_bindir}/%{name}
%{_datadir}/applications/%{name}.desktop
%{_userunitdir}/%{name}.service
%{_datadir}/icons/hicolor/256x256/apps/%{name}.png

%changelog
* Tue Aug 18 2026 Miraç Bedir - 2.1.1-3
- Add systemd daemon service, desktop entry, and application icon support

%global debug_package %{nil}
