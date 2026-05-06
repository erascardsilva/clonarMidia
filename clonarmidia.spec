Name:           clonarmidia
Version:        1.0.0
Release:        1%{?dist}
Summary:        Disk and Partition Cloning Tool

License:        MIT
URL:            https://github.com/erascardsilva/clonarMidia

Requires:       gtk3, webkit2gtk4.1

%description
Powerful tool for disk cloning and data recovery.

%install
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/usr/share/applications
mkdir -p %{buildroot}/usr/share/pixmaps
install -m 755 %{_sourcedir}/clonarmidia %{buildroot}/usr/bin/clonarmidia
install -m 644 %{_sourcedir}/clonarmidia.desktop %{buildroot}/usr/share/applications/clonarmidia.desktop
install -m 644 %{_sourcedir}/clonarmidia.png %{buildroot}/usr/share/pixmaps/clonarmidia.png

%files
/usr/bin/clonarmidia
/usr/share/applications/clonarmidia.desktop
/usr/share/pixmaps/clonarmidia.png

%changelog
* Tue May 05 2026 Erasmo Cardoso - 1.0.0-1
- Local build update
