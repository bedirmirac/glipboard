Name:           glipboard
Version:        2.1.1
Release:        1%{?dist}
Summary:        A terminal-based clipboard manager

License:        MIT
URL:            https://github.com/bedirmirac/glipboard
Source0:        %{url}/archive/v%{version}.tar.gz
Source1:        glipboard.desktop
Source2:        glipboard.service

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
# Binary dosyasını kur
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}

# .desktop dosyasını uygulama menüsüne kur
install -Dpm 0644 %{SOURCE1} %{buildroot}%{_datadir}/applications/%{name}.desktop

# Systemd kullanıcı servis dosyasını kur
install -Dpm 0644 %{SOURCE2} %{buildroot}%{_userunitdir}/%{name}.service

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

%changelog
* Tue Aug 18 2026 Miraç Bedir - 2.1.1-1
- Update to version 2.1.1
- Add systemd daemon service and desktop entry for TUI

%global debug_package %{nil}
